import logging
import os
import re
import tempfile
import base64
from typing import Any, Tuple, Dict, Union
from PIL import Image
import requests
from .base_parser import BaseParser
import pdfplumber
logger = logging.getLogger(__name__)


class PDFParser(BaseParser):
    """
    PDF Document Parser using either MinerU API (for scanned PDFs) or pymupdf4llm (for text-based PDFs).
    Returns parsed markdown text and a mapping of image URLs to PIL Image objects.
    """
    def _convert_table_to_markdown(self, table_data: list) -> str:
    
        if not table_data or not table_data[0]: return ""
        def clean_cell(cell):
            if cell is None: return ""
            return str(cell).replace("\n", " <br> ")
        try:
            markdown = ""
            header = [clean_cell(cell) for cell in table_data[0]]
            markdown += "| " + " | ".join(header) + " |\n"
            markdown += "| " + " | ".join(["---"] * len(header)) + " |\n"
            for row in table_data[1:]:
                if not row: continue
                body_row = [clean_cell(cell) for cell in row]
                if len(body_row) != len(header):
                    logger.warning(f"Skipping malformed table row: {body_row}")
                    continue
                markdown += "| " + " | ".join(body_row) + " |\n"
            return markdown
        except Exception as e:
            logger.error(f"Error converting table to markdown: {e}")
            return ""

    def parse_into_text(self, content: bytes) -> Union[str, Tuple[str, Dict[str, Any]]]:
        # Check if MinerU API is available
        api_healthy = False
        try:
            res = requests.get("http://host.docker.internal:7776/health")
            api_healthy = res.status_code == 200
        except requests.RequestException:
            logger.info("MinerU API is not reachable.")

        images_map: Dict[str, Image.Image] = {}

        # Helper to replace local image paths with uploaded URLs
        def replace_img(match):
                    prefix = match.group(1)
                    img_path = match.group(2)
                    suffix = match.group(3)
                    
                    if img_path.startswith(('http://', 'https://')):
                        return match.group(0)
                    
                    abs_img_path = f"{temp_dir}/{img_path}"
                    if not os.path.exists(abs_img_path):
                        logger.warning(f"警告：图片不存在，跳过: {abs_img_path}")
                    image = Image.open(abs_img_path).convert("RGBA")
                    image_url = self.upload_file(abs_img_path)
                    images_map[image_url] = image
                    return f"{prefix}{image_url}{suffix}"

        if api_healthy:
            logger.info(f"Parsing scanned PDF via MinerU API (size: {len(content)} bytes)")
            try:
                response = requests.post(
                    url="http://host.docker.internal:7776/file_parse",
                    data={
                        "return_md": True,
                        "return_images": True,
                        "lang_list": "ch",
                        "table_enable": True,
                        "formula_enable": True,
                        "parse_method": "auto",
                        "start_page_id": "0",
                        "end_page_id": "99999",
                        "backend": "pipeline",
                        "response_format_zip": False,
                        "return_middle_json": False,
                        "return_model_output": False,
                        "return_content_list": False,
                    },
                    files={"files": content},
                    timeout=1000,
                )
                response.raise_for_status()
                result = response.json()["results"]["files"]
                md_content = result["md_content"]
                images_b64 = result.get("images", {})

                with tempfile.TemporaryDirectory() as temp_dir:
                    images_dir = os.path.join(temp_dir, "images")
                    os.makedirs(images_dir, exist_ok=True)

                    for name, b64_str in images_b64.items():
                        # Strip data URI prefix if present
                        if b64_str.startswith("data:"):
                            b64_str = b64_str.split(",", 1)[1]
                        img_path = os.path.join(images_dir, name)
                        with open(img_path, "wb") as f:
                            f.write(base64.b64decode(b64_str))

                    # Replace image paths in markdown with uploaded URLs
                    img_pattern = r'(!\[.*?\]\()([^)\s]+)(\))'
                    text = re.sub(img_pattern,replace_img,md_content)

                return text, images_map

            except Exception as e:
                logger.error(f"MinerU parsing failed: {e}", exc_info=True)
                
        else:

            all_page_content = []
            temp_pdf = tempfile.NamedTemporaryFile(delete=False, suffix=".pdf")
            temp_pdf_path = temp_pdf.name
            logger.info(f"Parsing PDF with pdfplumber (size: {len(content)} bytes)")
            try:
                temp_pdf.write(content)
                temp_pdf.close()
                logger.info(f"PDF content written to temporary file: {temp_pdf_path}")

                with pdfplumber.open(temp_pdf_path) as pdf:
                    logger.info(f"PDF has {len(pdf.pages)} pages")
                    
                    for page_num, page in enumerate(pdf.pages):
                        page_content_parts = []
                        
                        # Try-fallback strategy for table detection
                        default_settings = { "vertical_strategy": "lines", "horizontal_strategy": "lines" }
                        found_tables = page.find_tables(default_settings)
                        if not found_tables:
                            logger.info(f"Page {page_num+1}: Default strategy found no tables. Trying fallback strategy.")
                            fallback_settings = { "vertical_strategy": "text", "horizontal_strategy": "lines" }
                            found_tables = page.find_tables(fallback_settings)

                        table_bboxes = [table.bbox for table in found_tables]
                        # Define a filter function that keeps objects NOT inside any table bbox.
                        def not_within_bboxes(obj):
                            """Check if an object is outside all table bounding boxes."""
                            for bbox in table_bboxes:
                                # Check if the object's vertical center is within a bbox
                                if bbox[1] <= (obj["top"] + obj["bottom"]) / 2 <= bbox[3]:
                                    return False # It's inside a table, so we DON'T keep it
                            return True # It's outside all tables, so we DO keep it

                        # that contains only the non-table text.
                        non_table_page = page.filter(not_within_bboxes)

                        # Now, extract text from this filtered page view.
                        text = non_table_page.extract_text(x_tolerance=2)
                        if text:
                            page_content_parts.append(text)
                
                        # Process and append the structured Markdown tables
                        if found_tables:
                            logger.info(f"Found {len(found_tables)} tables on page {page_num + 1}")
                            for table in found_tables:
                                markdown_table = self._convert_table_to_markdown(table.extract())
                                page_content_parts.append(f"\n\n{markdown_table}\n\n")
                        
                        
                        all_page_content.append("".join(page_content_parts))

                final_text = "\n\n--- Page Break ---\n\n".join(all_page_content)
                logger.info(f"PDF parsing complete. Extracted {len(final_text)} text chars.")
                
                return final_text
            except Exception as e:
                logger.error(f"Failed to parse PDF document: {str(e)}")
                return ""
            finally:
                # This block is GUARANTEED to execute, preventing resource leaks.
                if os.path.exists(temp_pdf_path):
                    try:
                        os.remove(temp_pdf_path)
                        logger.info(f"Temporary file cleaned up: {temp_pdf_path}")
                    except OSError as e:
                        logger.error(f"Error removing temporary file {temp_pdf_path}: {e}")

