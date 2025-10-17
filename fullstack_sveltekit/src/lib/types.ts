export type Picture = {
    aperture: string;
    camera: string;
    exposure: string;
    flash: string;
    focal: string;
    iso: string;
    lens: string;
    mode: string;
    dateString: string;
    path: string;
    landscape: boolean;
    notes: string;
};

export type Specie = {
    name: string;
    scientific_name: string;
    threat: string;
    info_page: string;
    description: string;
};

export type ReviewPicture = {
    path: string;
    name: string;
    hash: string;
    review_id: number;
    review_name: string;
    landscape: boolean;
    status: number;
    comment: string;
};

export type Album = {
    name: string;
    description: string;
};

export type AlbumWithPictures = {
    name: string;
    description: string;
    pictures: Array<Picture>;
}

export type Review = {
    id: number;
    name: string;
    password: string;
}