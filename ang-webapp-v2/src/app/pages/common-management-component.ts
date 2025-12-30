import { inject, signal } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";

export class CommonLayoutComponent<T extends { nimb_id: string; audit: { version: number } }> {

    viewMode = signal<'list' | 'editor'>('list');
    router = inject(Router);
    activatedRoute = inject(ActivatedRoute);

    constructor() {
        this.activatedRoute.queryParams.subscribe(params => {
            const editId = params['nimb_id'];
            const editVersion = params['version'];
            if (editId && editVersion) {
                this.viewMode.set('editor');
            } else {
                this.viewMode.set('list');
            }
        });
    }

    setViewMode(mode: 'list' | 'editor') {
        this.viewMode.set(mode);
    }

    handleEdit(data: T) {
        this.router.navigate([], {
            queryParams: {
                nimb_id: data.nimb_id,
                version: data.audit.version
            }
        });
    }
}