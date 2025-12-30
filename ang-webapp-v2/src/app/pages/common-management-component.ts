import { inject, signal } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";

export class CommonManagementComponent {

    viewMode = signal<'list' | 'editor'>('list');
    router = inject(Router);
    activatedRoute = inject(ActivatedRoute)


    setViewMode(mode: 'list' | 'editor') {
        this.viewMode.set(mode);
    }
}