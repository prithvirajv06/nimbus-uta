export interface Audit {
    created_at: Date;
    created_by: string;
    modified_at: Date;
    modified_by: string;
    is_prod_candidate: boolean;
    status: string
    active: string;
    is_archived: boolean;
    restore_archive: boolean
    version: number;
    minor_version: number;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    pageSize: number;
}

export interface SelectOption {
    label: string;
    value: string | number;
}

export interface ApiResponse<T> {
    status: string;
    message: string;
    data: T;
}