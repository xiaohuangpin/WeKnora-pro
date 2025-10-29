> **WeKnora‑pro** is a derivative version based on the original **WeKnora**, primarily enhanced for advanced document parsing capabilities.  
> Key improvements include:
> * Support for scanned documents via the **Mineru‑API** (with automatic CPU/GPU optimization) for OCR and table extraction, fully compatible with WeKnora’s multimodal pipeline.
> * Increased maximum document size limit to **300 MB**.

<p align="center">
  <picture>
    <img src="./docs/images/logo.png" alt="WeKnora Logo" height="120"/>
  </picture>
</p>

---

## 📌 Project Overview

- **Enhanced Parsing**: Integrated a new PDF parser powered by the Mineru backend.
- **Large File Support**: Supports single files up to **300 MB**.
- **WeChat Ecosystem Compatibility**: Seamlessly integrates into WeChat Official Accounts, Mini Programs, and other WeChat scenarios.

---

## 🚀 Quick Start

> All instructions below assume a Linux/macOS environment (Windows is also supported). Please ensure the following tools are installed:

| Tool | Official Website |
|------|------------------|
| Docker | https://www.docker.com/ |
| Docker Compose | https://docs.docker.com/compose/ |
| Git | https://git-scm.com/ |

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/Tencent/WeKnora.git
cd WeKnora
```

### 2️⃣ Configure Environment Variables

```bash
cp .env.example .env      # Copy the example configuration file
# Edit .env to configure database, Redis, OpenAI, etc.
# Note: For storage backend, MinIO is recommended; the 'local' option may prevent multi-model enhancement from functioning properly.
```

> Detailed descriptions for all variables are provided as comments in `.env.example`.

### 3️⃣ Build Images and Start Services

```bash
# ① Build Docker images (including Ollama and backend containers)
./scripts/build_images.sh

# ② Start all services (by default, does not pull latest images; remove --no-pull to update)
./scripts/start_all.sh --no-pull
```

> **⏰** First-time startup may take **5–10 minutes** due to image downloads and initialization.

### 4️⃣ Launch the Mineru Service (Advanced Document Parser)

```bash
# Create a Python virtual environment
conda create -n mineru python=3.10
conda activate mineru

# Install dependencies
pip install uv -i https://pypi.tuna.tsinghua.edu.cn/simple     # If network issues occur, use plain pip
uv pip install -U "mineru[core]"  # Launch the API server
# For faster parsing on GPUs with >16GB VRAM, consider using "mineru[all]"

cd wed_api
python web_api.py
```

> **🛑** After document parsing is complete, press `Ctrl+C` to stop the service and release GPU memory. This will **not** affect the Q&A functionality.

### 5️⃣ Stop All Services

```bash
./scripts/start_all.sh --stop   # or
make stop-all
```

---

## 🌐 Service Endpoints

| Service | URL |
|--------|-----|
| Web UI | `http://localhost` |
| Backend API | `http://localhost:8080` |
| Jaeger Tracing | `http://localhost:16686` |

> If deployed on a remote server, replace `localhost` with the appropriate host IP address.

---

## 🔧 Contribution Guidelines

1. Fork this repository.  
2. Clone your fork locally:

   ```bash
   git clone git@github.com:<your-username>/WeKnora-pro.git
   cd WeKnora-pro
   ```

3. Create and switch to a new feature branch:

   ```bash
   git checkout -b feature/<brief-description>
   ```

4. Before submitting, ensure that:

   * Code adheres to **PEP8** (Python) and **Prettier** (JavaScript) standards.
   * Unit/integration tests are added or updated, with **≥80% coverage**.
   * Documentation and README files are updated accordingly.

5. Push your branch and open a Pull Request:

   ```bash
   git push origin feature/<brief-description>
   ```

6. Await review from maintainers.

---

## 🙏 Acknowledgements

This project builds upon the following open-source components:

- [WeKnora](https://github.com/Tencent/WeKnora)
- [MinerU](https://github.com/opendatalab/MinerU)

---

## 📜 License & Usage Restrictions

1. **AGPL‑v3 License**:  
   * All derivative works must also be licensed under AGPL‑v3.  
   * If the software is offered as a network service, users must be able to request and receive the corresponding source code (per Section 13 of AGPL‑v3).

2. **Commercial Use**:  
   * Commercial usage—including SaaS offerings and internal enterprise deployments—is permitted.  
   * Even if the original code is unmodified, full source code disclosure and compliance with AGPL‑v3 requirements are mandatory.  
   * For proprietary (closed-source) commercial use, **written authorization** from all copyright holders is required.

3. **Disclaimer**:  
   This project is provided **“as is”**, without any warranties of any kind. Users are solely responsible for assessing legal and compliance risks. Consult a qualified legal professional when necessary.