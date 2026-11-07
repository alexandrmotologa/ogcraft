import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Template, TemplateSummary } from './lib/types';
import {
  listTemplates,
  getTemplate,
  saveTemplate,
  deleteTemplate,
  renderLivePreview,
} from './lib/api';
import { CanvasPreview } from './components/CanvasPreview';
import { TemplateSelector } from './components/TemplateSelector';
import { ElementInspector } from './components/ElementInspector';
import { ParamTester } from './components/ParamTester';
import { CodeExporter } from './components/CodeExporter';
import { Cpu, Save, ShieldCheck } from 'lucide-react';

export const App: React.FC = () => {
  const [templateSummaries, setTemplateSummaries] = useState<TemplateSummary[]>([]);
  const [activeTemplateId, setActiveTemplateId] = useState<string>('dark_tech');
  const [activeTemplate, setActiveTemplate] = useState<Template | null>(null);
  const [selectedElementIndex, setSelectedElementIndex] = useState<number | null>(null);

  const [params, setParams] = useState<Record<string, string>>({});
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [renderTime, setRenderTime] = useState<string>('');
  const [isSaving, setIsSaving] = useState<boolean>(false);
  const [saveStatus, setSaveStatus] = useState<string | null>(null);

  const debounceTimerRef = useRef<number | null>(null);

  // 1. Fetch template summaries
  const loadTemplatesList = useCallback(async () => {
    try {
      const list = await listTemplates();
      setTemplateSummaries(list);
      if (list.length > 0 && !list.some((t) => t.id === activeTemplateId)) {
        setActiveTemplateId(list[0].id);
      }
    } catch (e) {
      console.error('Failed to list templates', e);
    }
  }, [activeTemplateId]);

  useEffect(() => {
    loadTemplatesList();
  }, [loadTemplatesList]);

  // 2. Fetch active template schema
  useEffect(() => {
    if (!activeTemplateId) return;

    let cancelled = false;
    getTemplate(activeTemplateId)
      .then((tmpl) => {
        if (cancelled) return;
        setActiveTemplate(tmpl);
        setSelectedElementIndex(null);

        // Pre-fill params with element defaults
        const initialParams: Record<string, string> = {};
        for (const el of tmpl.elements) {
          if (el.field && el.default) {
            initialParams[el.field] = el.default;
          }
        }
        setParams(initialParams);
      })
      .catch((err) => console.error(err));

    return () => {
      cancelled = true;
    };
  }, [activeTemplateId]);

  // 3. Render card preview on template or params change (debounced 150ms)
  const triggerRender = useCallback(
    (tmpl: Template, currentParams: Record<string, string>) => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }

      debounceTimerRef.current = window.setTimeout(async () => {
        setLoading(true);
        const startTime = performance.now();
        try {
          const url = await renderLivePreview(tmpl, currentParams);
          setPreviewUrl(url);
          const duration = (performance.now() - startTime).toFixed(1);
          setRenderTime(`${duration}ms`);
        } catch (e) {
          console.error('Failed to render live preview', e);
        } finally {
          setLoading(false);
        }
      }, 150);
    },
    []
  );

  useEffect(() => {
    if (!activeTemplate) return;
    triggerRender(activeTemplate, params);
  }, [activeTemplate, params, triggerRender]);

  const handleSelectTemplate = (id: string) => {
    setActiveTemplateId(id);
  };

  const handleCreateNewTemplate = () => {
    const newId = `custom_${Date.now().toString().slice(-4)}`;
    const newTemplate: Template = {
      id: newId,
      name: `New Custom Card`,
      description: 'Custom social preview card',
      width: 1200,
      height: 630,
      background: {
        type: 'gradient',
        direction: 'to-bottom-right',
        stops: [
          { offset: 0.0, color: '#090d16' },
          { offset: 1.0, color: '#161e2e' },
        ],
      },
      elements: [
        {
          type: 'badge',
          field: 'tag',
          default: 'CUSTOM',
          font: 'JetBrainsMono-Medium',
          font_size: 18,
          color: '#ffffff',
          background_color: '#2563eb',
          x: 80,
          y: 80,
          padding_x: 14,
          padding_y: 6,
          border_radius: 6,
        },
        {
          type: 'text',
          field: 'title',
          default: 'Your Custom Headline Here',
          font: 'Inter-Bold',
          font_size: 54,
          color: '#f8fafc',
          x: 80,
          y: 160,
          max_width: 1040,
          line_height: 1.25,
          max_lines: 3,
          fit_text: true,
        },
      ],
    };

    setActiveTemplate(newTemplate);
    setActiveTemplateId(newId);
  };

  const handleSaveCurrentTemplate = async () => {
    if (!activeTemplate) return;
    setIsSaving(true);
    try {
      await saveTemplate(activeTemplate);
      setSaveStatus('Template saved successfully!');
      setTimeout(() => setSaveStatus(null), 2500);
      loadTemplatesList();
    } catch (e: any) {
      alert(e.message || 'Failed to save template');
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeleteCustomTemplate = async (id: string) => {
    try {
      await deleteTemplate(id);
      loadTemplatesList();
      if (activeTemplateId === id) {
        setActiveTemplateId('dark_tech');
      }
    } catch (e: any) {
      alert(e.message || 'Failed to delete template');
    }
  };

  const handleParamChange = (field: string, value: string) => {
    setParams((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  return (
    <div className="flex flex-col h-screen w-screen overflow-hidden bg-[#070a10] text-slate-100 font-sans">
      {/* Top Navbar */}
      <header className="h-14 border-b border-[#1a2333] bg-[#0c121e] px-5 flex items-center justify-between flex-shrink-0 z-20">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-sky-500/15 border border-sky-500/30 flex items-center justify-center font-mono font-bold text-sky-400">
              OG
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold text-sm tracking-tight text-white">ogcraft</span>
                <span className="text-[10px] px-1.5 py-0.2 rounded bg-sky-500/20 text-sky-400 font-mono font-semibold">
                  v1.0
                </span>
              </div>
              <p className="text-[10px] text-slate-400 -mt-0.5">Real-Time OpenGraph Studio</p>
            </div>
          </div>

          <div className="h-5 w-px bg-slate-800 mx-1 hidden sm:block" />

          <div className="hidden md:flex items-center gap-4 text-xs font-mono text-slate-400">
            <span className="flex items-center gap-1.5">
              <Cpu size={13} className="text-emerald-400" />
              SIMD Go 2D
            </span>
            <span className="flex items-center gap-1.5">
              <ShieldCheck size={13} className="text-sky-400" />
              SSRF & HMAC Protected
            </span>
          </div>
        </div>

        {/* Right Action Buttons */}
        <div className="flex items-center gap-3">
          {saveStatus && (
            <span className="text-xs text-emerald-400 font-mono transition animate-fade-in">
              ✓ {saveStatus}
            </span>
          )}

          <button
            onClick={handleSaveCurrentTemplate}
            disabled={isSaving || !activeTemplate}
            className="flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium border border-slate-700 transition"
          >
            <Save size={13} />
            {isSaving ? 'Saving...' : 'Save Template'}
          </button>
        </div>
      </header>

      {/* Main Studio Body */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left: Template Selector */}
        <TemplateSelector
          templates={templateSummaries}
          activeId={activeTemplateId}
          onSelect={handleSelectTemplate}
          onCreateNew={handleCreateNewTemplate}
          onDelete={handleDeleteCustomTemplate}
        />

        {/* Center: Canvas Viewport & Parameter Controls */}
        {activeTemplate ? (
          <div className="flex-1 flex flex-col h-full overflow-hidden">
            <CanvasPreview
              template={activeTemplate}
              previewUrl={previewUrl}
              loading={loading}
              renderTime={renderTime}
              selectedElementIndex={selectedElementIndex}
              onSelectElement={setSelectedElementIndex}
              onUpdateElement={(idx, updated) => {
                const nextElements = [...activeTemplate.elements];
                nextElements[idx] = updated;
                setActiveTemplate({
                  ...activeTemplate,
                  elements: nextElements,
                });
              }}
              onRefresh={() => triggerRender(activeTemplate, params)}
            />

            {/* Bottom Controls: URL Parameter Simulator & Export */}
            <div className="flex-shrink-0 flex flex-col max-h-72 overflow-y-auto">
              <ParamTester
                template={activeTemplate}
                params={params}
                onChangeParam={handleParamChange}
                onResetParams={() => {
                  const initial: Record<string, string> = {};
                  for (const el of activeTemplate.elements) {
                    if (el.field && el.default) {
                      initial[el.field] = el.default;
                    }
                  }
                  setParams(initial);
                }}
              />
              <CodeExporter template={activeTemplate} params={params} />
            </div>
          </div>
        ) : (
          <div className="flex-1 flex items-center justify-center text-slate-500 font-mono text-sm">
            Loading template...
          </div>
        )}

        {/* Right: Element / Canvas Inspector */}
        {activeTemplate && (
          <ElementInspector
            template={activeTemplate}
            selectedElementIndex={selectedElementIndex}
            onUpdateTemplate={setActiveTemplate}
            onSelectElement={setSelectedElementIndex}
          />
        )}
      </div>
    </div>
  );
};
