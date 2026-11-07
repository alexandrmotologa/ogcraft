import React, { useRef, useState, useEffect } from 'react';
import { ElementSpec, Template } from '../lib/types';
import { Maximize2, ZoomIn, ZoomOut, RefreshCw, Download } from 'lucide-react';

interface CanvasPreviewProps {
  template: Template;
  previewUrl: string | null;
  loading: boolean;
  renderTime: string;
  selectedElementIndex: number | null;
  onSelectElement: (index: number | null) => void;
  onUpdateElement: (index: number, updated: ElementSpec) => void;
  onRefresh: () => void;
}

export const CanvasPreview: React.FC<CanvasPreviewProps> = ({
  template,
  previewUrl,
  loading,
  renderTime,
  selectedElementIndex,
  onSelectElement,
  onUpdateElement,
  onRefresh,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [scale, setScale] = useState<number>(0.7);
  const [isDragging, setIsDragging] = useState<boolean>(false);
  const [dragIndex, setDragIndex] = useState<number | null>(null);
  const [dragOffset, setDragOffset] = useState<{ x: number; y: number }>({ x: 0, y: 0 });

  // Auto-fit scale on mount and resize
  useEffect(() => {
    const updateScale = () => {
      if (!containerRef.current) return;
      const containerWidth = containerRef.current.clientWidth - 48;
      const targetScale = Math.min(1.0, Math.max(0.35, containerWidth / 1200));
      setScale(targetScale);
    };

    updateScale();
    window.addEventListener('resize', updateScale);
    return () => window.removeEventListener('resize', updateScale);
  }, []);

  const handleMouseDown = (e: React.MouseEvent, index: number) => {
    e.stopPropagation();
    onSelectElement(index);
    setIsDragging(true);
    setDragIndex(index);

    const el = template.elements[index];
    setDragOffset({
      x: e.clientX / scale - el.x,
      y: e.clientY / scale - el.y,
    });
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging || dragIndex === null) return;
    const currentEl = template.elements[dragIndex];

    const rawX = e.clientX / scale - dragOffset.x;
    const rawY = e.clientY / scale - dragOffset.y;

    // Snap to 10px grid
    const newX = Math.max(0, Math.min(1200 - 50, Math.round(rawX / 10) * 10));
    const newY = Math.max(0, Math.min(630 - 30, Math.round(rawY / 10) * 10));

    if (newX !== currentEl.x || newY !== currentEl.y) {
      onUpdateElement(dragIndex, {
        ...currentEl,
        x: newX,
        y: newY,
      });
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
    setDragIndex(null);
  };

  const handleDownload = () => {
    if (!previewUrl) return;
    const a = document.createElement('a');
    a.href = previewUrl;
    a.download = `${template.id || 'card'}.png`;
    a.click();
  };

  return (
    <div className="flex flex-col flex-1 h-full bg-[#070a10] border-r border-[#1a2333] overflow-hidden">
      {/* Canvas Top Bar */}
      <div className="flex items-center justify-between px-6 py-3 border-b border-[#1a2333] bg-[#0c121e]">
        <div className="flex items-center gap-3">
          <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">Preview</span>
          <span className="text-xs px-2 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">
            1200 × 630 px
          </span>
          {renderTime && (
            <span className="text-xs px-2 py-0.5 rounded bg-emerald-950/70 text-emerald-400 border border-emerald-800/40 font-mono">
              ⚡ {renderTime}
            </span>
          )}
        </div>

        <div className="flex items-center gap-2">
          <div className="flex items-center bg-slate-900 border border-slate-800 rounded-lg p-0.5">
            <button
              onClick={() => setScale(s => Math.max(0.3, s - 0.1))}
              className="p-1.5 hover:bg-slate-800 text-slate-400 hover:text-slate-200 rounded"
              title="Zoom out"
            >
              <ZoomOut size={14} />
            </button>
            <span className="text-xs text-slate-300 font-mono px-2 min-w-12 text-center">
              {Math.round(scale * 100)}%
            </span>
            <button
              onClick={() => setScale(s => Math.min(1.0, s + 0.1))}
              className="p-1.5 hover:bg-slate-800 text-slate-400 hover:text-slate-200 rounded"
              title="Zoom in"
            >
              <ZoomIn size={14} />
            </button>
            <button
              onClick={() => setScale(0.7)}
              className="p-1.5 hover:bg-slate-800 text-slate-400 hover:text-slate-200 rounded ml-0.5"
              title="Reset Zoom"
            >
              <Maximize2 size={14} />
            </button>
          </div>

          <button
            onClick={onRefresh}
            disabled={loading}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-slate-200 bg-slate-800 hover:bg-slate-700 rounded-lg transition"
          >
            <RefreshCw size={13} className={loading ? 'animate-spin' : ''} />
            Refresh
          </button>

          <button
            onClick={handleDownload}
            disabled={!previewUrl}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-slate-900 bg-sky-400 hover:bg-sky-300 rounded-lg transition font-semibold"
          >
            <Download size={13} />
            Export PNG
          </button>
        </div>
      </div>

      {/* Canvas Viewport */}
      <div
        ref={containerRef}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onClick={() => onSelectElement(null)}
        className="flex-1 flex items-center justify-center p-6 overflow-auto bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:16px_16px] relative"
      >
        <div
          style={{
            width: 1200 * scale,
            height: 630 * scale,
            transformOrigin: 'center center',
          }}
          className="relative rounded-lg shadow-2xl shadow-black/80 overflow-hidden ring-1 ring-slate-800 bg-black flex-shrink-0"
        >
          {/* Rendered Background / Image */}
          {previewUrl ? (
            <img
              src={previewUrl}
              alt="OG Card Preview"
              className="w-full h-full object-contain pointer-events-none select-none"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-slate-500 font-mono text-sm">
              Generating Preview...
            </div>
          )}

          {/* Interactive Element Overlays */}
          {template.elements.map((el, idx) => {
            const isSelected = selectedElementIndex === idx;
            const w = el.max_width || el.width || (el.size ? el.size : 200);
            const h = el.height || (el.size ? el.size : (el.font_size ? el.font_size * 1.5 : 40));

            return (
              <div
                key={idx}
                onMouseDown={(e) => handleMouseDown(e, idx)}
                style={{
                  left: el.x * scale,
                  top: el.y * scale,
                  width: Math.max(40, w * scale),
                  minHeight: Math.max(24, h * scale),
                }}
                className={`absolute cursor-move transition-all select-none group ${
                  isSelected
                    ? 'ring-2 ring-sky-400 bg-sky-500/10 rounded'
                    : 'hover:ring-1 hover:ring-sky-400/50 hover:bg-white/5 rounded'
                }`}
              >
                {/* Element Badge */}
                {isSelected && (
                  <div className="absolute -top-6 left-0 px-1.5 py-0.5 rounded bg-sky-500 text-[10px] font-mono text-slate-950 font-bold tracking-wider uppercase pointer-events-none whitespace-nowrap shadow">
                    {el.type}: {el.field || el.type} ({el.x}, {el.y})
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
