import { get, post, put, del, postUpload, getDown } from "../../utils/request";

// 知识库管理 API（列表、创建、获取、更新、删除、复制）
export function listKnowledgeBases() {
  return get(`/api/v1/knowledge-bases`);
}

export function createKnowledgeBase(data: { name: string; description?: string; chunking_config?: any }) {
  return post(`/api/v1/knowledge-bases`, data);
}

export function getKnowledgeBaseById(id: string) {
  return get(`/api/v1/knowledge-bases/${id}`);
}

export function updateKnowledgeBase(id: string, data: { name: string; description?: string; config: any }) {
  return put(`/api/v1/knowledge-bases/${id}` , data);
}

export function deleteKnowledgeBase(id: string) {
  return del(`/api/v1/knowledge-bases/${id}`);
}

export function copyKnowledgeBase(data: { source_id: string; target_id?: string }) {
  return post(`/api/v1/knowledge-bases/copy`, data);
}

// 知识文件 API（基于具体知识库）
export function uploadKnowledgeFile(kbId: string, data = {}) {
  return postUpload(`/api/v1/knowledge-bases/${kbId}/knowledge/file`, data);
}

export function listKnowledgeFiles(kbId: string, { page, page_size }: { page: number; page_size: number }) {
  return get(`/api/v1/knowledge-bases/${kbId}/knowledge?page=${page}&page_size=${page_size}`);
}

export function getKnowledgeDetails(id: string) {
  return get(`/api/v1/knowledge/${id}`);
}

export function delKnowledgeDetails(id: string) {
  return del(`/api/v1/knowledge/${id}`);
}

export function downKnowledgeDetails(id: string) {
  return getDown(`/api/v1/knowledge/${id}/download`);
}

export function batchQueryKnowledge(idsQueryString: string) {
  return get(`/api/v1/knowledge/batch?${idsQueryString}`);
}

export function getKnowledgeDetailsCon(id: string, page: number) {
  return get(`/api/v1/chunks/${id}?page=${page}&page_size=25`);
}