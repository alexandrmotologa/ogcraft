import { Template, TemplateSummary } from './types';

const API_BASE = '';

export async function listTemplates(): Promise<TemplateSummary[]> {
  const res = await fetch(`${API_BASE}/v1/templates`);
  if (!res.ok) throw new Error('Failed to fetch templates');
  return res.json();
}

export async function getTemplate(id: string): Promise<Template> {
  const res = await fetch(`${API_BASE}/v1/templates/${id}`);
  if (!res.ok) throw new Error(`Failed to fetch template ${id}`);
  return res.json();
}

export async function saveTemplate(template: Template): Promise<void> {
  const res = await fetch(`${API_BASE}/v1/templates`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(template),
  });
  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || 'Failed to save template');
  }
}

export async function deleteTemplate(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/v1/templates/${id}`, {
    method: 'DELETE',
  });
  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || 'Failed to delete template');
  }
}

export async function renderLivePreview(template: Template, params: Record<string, string>): Promise<string> {
  const res = await fetch(`${API_BASE}/v1/render`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ template, params }),
  });
  if (!res.ok) {
    throw new Error('Failed to render preview');
  }
  const blob = await res.blob();
  return URL.createObjectURL(blob);
}

export function buildCardUrl(templateId: string, params: Record<string, string>): string {
  const url = new URL(`${window.location.origin}/v1/cards/${templateId}.png`);
  for (const [key, val] of Object.entries(params)) {
    if (val && val.trim() !== '') {
      url.searchParams.set(key, val);
    }
  }
  return url.toString();
}
