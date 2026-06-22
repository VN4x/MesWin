import { pb } from './pb';

export type MeasuringOrder = {
	id: string;
	status: string;
	notes: string;
	customer: string;
	organization: string;
	created: string;
	updated: string;
	expand?: {
		customer?: Customer;
	};
};

export type Customer = {
	id: string;
	name: string;
	address: string;
	contact: string;
	email: string;
	phone: string;
	notes: string;
	organization: string;
};

export type WindowItem = {
	id: string;
	order: string;
	label: string;
	window_type: string;
	status: string;
	inner_photos: string[];
	outer_photos: string[];
	cv_result?: Record<string, unknown>;
	overrides?: Record<string, unknown>;
	confidence?: number;
};

export async function listOrders(): Promise<MeasuringOrder[]> {
	return pb.collection('measuring_orders').getFullList<MeasuringOrder>({
		sort: '-created',
		expand: 'customer'
	});
}

export async function getOrder(id: string): Promise<MeasuringOrder> {
	return pb.collection('measuring_orders').getOne<MeasuringOrder>(id, { expand: 'customer' });
}

export async function createCustomer(data: Partial<Customer>) {
	return pb.collection('customers').create(data);
}

export async function createOrder(data: {
	organization: string;
	customer: string;
	status?: string;
	notes?: string;
}) {
	return pb.collection('measuring_orders').create({
		...data,
		status: data.status ?? 'draft'
	});
}

export async function listWindows(orderId: string): Promise<WindowItem[]> {
	return pb.collection('window_items').getFullList<WindowItem>({
		filter: `order = "${orderId}"`,
		sort: 'created'
	});
}

export async function addWindow(orderId: string, windowType: string, label?: string) {
	return pb.collection('window_items').create({
		order: orderId,
		window_type: windowType,
		label: label ?? '',
		status: 'pending'
	});
}

export async function processWindow(windowId: string) {
	return pb.send(`/api/meswin/windows/${windowId}/process`, { method: 'POST' });
}

export async function exportOrderJson(orderId: string) {
	return pb.send(`/api/meswin/orders/${orderId}/export.json`, { method: 'GET' });
}
