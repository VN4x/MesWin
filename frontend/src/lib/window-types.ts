export type WindowTypeSlug =
	| 'single_fixed'
	| 'single_tilt_turn'
	| 'single_side_hung'
	| 'single_tilt_only'
	| 'double_fixed_operable'
	| 'double_two_tilt_turn'
	| 'double_french_casement'
	| 'triple_fixed_center_operable_fixed'
	| 'triple_operable_fixed_operable'
	| 'triple_mixed'
	| 'multi_4plus'
	| 'multi_transom'
	| 'multi_mullion';

export const WINDOW_TYPES: { slug: WindowTypeSlug; label: string; group: string }[] = [
	{ slug: 'single_fixed', label: 'Single — Fixed', group: 'Single' },
	{ slug: 'single_tilt_turn', label: 'Single — Tilt & Turn', group: 'Single' },
	{ slug: 'single_side_hung', label: 'Single — Side-Hung Casement', group: 'Single' },
	{ slug: 'single_tilt_only', label: 'Single — Tilt-Only', group: 'Single' },
	{ slug: 'double_fixed_operable', label: 'Double — Fixed + Operable', group: 'Double' },
	{ slug: 'double_two_tilt_turn', label: 'Double — Two Tilt & Turn', group: 'Double' },
	{ slug: 'double_french_casement', label: 'Double — French Casement', group: 'Double' },
	{ slug: 'triple_fixed_center_operable_fixed', label: 'Triple — Fixed + Center + Fixed', group: 'Triple' },
	{ slug: 'triple_operable_fixed_operable', label: 'Triple — Operable + Fixed + Operable', group: 'Triple' },
	{ slug: 'triple_mixed', label: 'Triple — Mixed', group: 'Triple' },
	{ slug: 'multi_4plus', label: 'Multi — 4+ Units', group: 'Multi' },
	{ slug: 'multi_transom', label: 'Multi — Transom', group: 'Multi' },
	{ slug: 'multi_mullion', label: 'Multi — With Mullions', group: 'Multi' }
];

export function windowTypeLabel(slug: string): string {
	return WINDOW_TYPES.find((t) => t.slug === slug)?.label ?? slug;
}
