

export interface ManagementContract {
    loadtableView(): void;
    loadtableArchivedView(): void;
    restore(item: any): void;
    archive(item: any): void;
    editDetails(item: any): void;
    setService(service: any): void;
}