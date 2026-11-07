export interface GradientStop {
  offset: number;
  color: string;
}

export interface BackgroundSpec {
  type: 'solid' | 'gradient' | 'radial';
  color?: string;
  direction?: 'to-bottom' | 'to-right' | 'to-bottom-right' | 'to-bottom-left';
  stops?: GradientStop[];
  center_x?: number;
  center_y?: number;
  radius?: number;
  center_color?: string;
  outer_color?: string;
  pattern?: 'none' | 'dots' | 'grid' | 'lines';
  pattern_color?: string;
  pattern_spacing?: number;
}

export interface ElementSpec {
  type: 'text' | 'badge' | 'avatar' | 'icon' | 'divider' | 'container' | 'code' | 'metric';
  field?: string;
  default?: string;
  font?: string;
  font_size?: number;
  color?: string;
  background_color?: string;
  border_color?: string;
  border_width?: number;
  border_radius?: number;
  x: number;
  y: number;
  width?: number;
  height?: number;
  max_width?: number;
  line_height?: number;
  max_lines?: number;
  fit_text?: boolean;
  padding_x?: number;
  padding_y?: number;
  icon?: string;
  size?: number;
  shape?: 'circle' | 'rounded' | 'square';

  // Code snippet properties
  language?: string;
  show_controls?: boolean;
  line_numbers?: boolean;

  // Metric properties
  value?: string;
  subtext?: string;
  progress?: number;
}

export interface Template {
  id: string;
  name: string;
  description?: string;
  width: number;
  height: number;
  background: BackgroundSpec;
  elements: ElementSpec[];
  builtin?: boolean;
}

export interface TemplateSummary {
  id: string;
  name: string;
  description?: string;
  builtin: boolean;
  elements_count: number;
}
