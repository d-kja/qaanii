type Image = string;

export interface Page {
	order: number;
	image: Image;
	imageUrl: string;
	imageType: string;
}

export interface Chapter {
	slug: string;
	date?: string;
	title?: string;

	pages?: Record<string, Page>;
}

export interface Manga {
	image: Image;
	image_url: string;
	image_type: string;

	name: string;
	slug: string;
	description?: string;

	tags?: string[] | null;
	url: string;

	chapters?: Chapter[];
}
