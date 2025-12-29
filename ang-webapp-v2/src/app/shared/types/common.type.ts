export interface Audit {
    created_at?: string;
    created_by ?: string;
    updated_at?: string;
    updated_by ?: string;
    is_prod_candidate?: boolean;
    active?: string;
    is_archived?: boolean;
    version?: number;
    minor_version?: number;
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