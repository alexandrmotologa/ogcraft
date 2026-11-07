import React, { useState } from 'react';
import { Template } from '../lib/types';
import { X, Upload, FileCode } from 'lucide-react';

interface TemplateImportModalProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (template: Template) => void;
}

export const TemplateImportModal: React.FC<TemplateImportModalProps> = ({
  isOpen,
  onClose,
  onImport,
}) => {
  const [jsonText, setJsonText] = useState('');
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleImport = () => {
    try {
      setError(null);
      const parsed = JSON.parse(jsonText);
      if (!parsed.id || !parsed.name || !Array.isArray(parsed.elements)) {
        throw new Error('Invalid template schema: must contain id, name, and elements array');
      }
      onImport(parsed);
      onClose();
    } catch (e: any) {
      setError(e.message || 'Failed to parse template JSON');
    }
  };

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (event) => {
      setJsonText((event.target?.result as string) || '');
    };
    reader.readAsText(file);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/75 backdrop-blur-sm p-4">
      <div className="bg-[#0c121e] border border-slate-800 rounded-2xl w-full max-w-xl p-6 shadow-2xl flex flex-col gap-4 animate-scale-in">
        <div className="flex items-center justify-between border-b border-slate-800 pb-3">
          <div className="flex items-center gap-2">
            <FileCode size={18} className="text-sky-400" />
            <h3 className="text-sm font-bold text-slate-100">Import Template JSON</h3>
          </div>
          <button
            onClick={onClose}
            className="p-1 hover:bg-slate-800 text-slate-400 hover:text-slate-200 rounded transition"
          >
            <X size={16} />
          </button>
        </div>

        <p className="text-xs text-slate-400">
          Paste valid template JSON or upload a <code>.json</code> file to immediately load it into the designer:
        </p>

        <textarea
          rows={10}
          value={jsonText}
          onChange={(e) => setJsonText(e.target.value)}
          placeholder={`{\n  "id": "my_template",\n  "name": "My Custom Template",\n  "width": 1200,\n  "height": 630,\n  "elements": []\n}`}
          className="w-full bg-[#070a10] border border-slate-800 rounded-lg p-3 text-xs font-mono text-slate-200 focus:outline-none focus:border-sky-500/50"
        />

        {error && <div className="text-xs text-rose-400 font-mono">⚠️ {error}</div>}

        <div className="flex items-center justify-between pt-2">
          <label className="flex items-center gap-1.5 px-3 py-1.5 rounded bg-slate-900 border border-slate-800 hover:border-slate-700 text-xs text-slate-300 cursor-pointer transition">
            <Upload size={14} />
            Choose JSON File
            <input type="file" accept=".json" onChange={handleFileUpload} className="hidden" />
          </label>

          <div className="flex items-center gap-2">
            <button
              onClick={onClose}
              className="px-4 py-1.5 rounded text-xs font-medium text-slate-400 hover:text-slate-200 transition"
            >
              Cancel
            </button>
            <button
              onClick={handleImport}
              disabled={!jsonText.trim()}
              className="px-4 py-1.5 rounded bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs transition disabled:opacity-50"
            >
              Load Template
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
