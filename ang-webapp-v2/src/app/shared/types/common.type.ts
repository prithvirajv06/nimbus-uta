export interface Audit {
    created_at: string;
    created_by: string;
    modified_at: string;
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
    success: boolean;
    message: string;
    data: T;
}