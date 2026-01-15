import { Component, inject, OnInit, WritableSignal } from "@angular/core";
import { FieldTree } from "@angular/forms/signals";
import { ActivatedRoute, Router } from "@angular/router";
import { NotificationService } from "../../services/notification.service";
import { ApiResponse } from "../../types/common.type";
import { RulesCommons } from "../../util-fulctions/options";

export interface CommonEditorContract<T> {
    setService(): void;
    setFormModel(): void;
    setFormGroup(): void;
    saveDetails(): void;
    getDetails(): void;
}

/**
 * Generic base component for editing entities of type `T`.
 * 
 * This component provides common logic for loading, saving, and cancelling edits,
 * using Angular's dependency injection for routing and notifications.
 * 
 * @typeParam T - The type of the entity being edited.
 * 
 * ## Usage
 * Extend this class and implement the abstract methods:
 * - `setService()`
 * - `setFormModel()`
 * - `setFormGroup()`
 * 
 * These methods must be overridden in the subclass to provide specific service,
 * form model, and form group logic for the entity type.
 * 
 * ## Methods
 * - `getDetails()` - Loads entity details based on query parameters.
 * - `saveDetails()` - Saves the current form data using the service.
 * - `cancelEdit()` - Resets the form and navigates away.
 * 
 * ## Note
 * To ensure subclasses must override methods, mark them as `abstract` in the base class.
 * For example:
 * ```typescript
 * abstract setService(): void;
 * ```
 * This will enforce implementation in derived classes.
 */
@Component({
    selector: 'app-common-editor',
    template: ''
})
export class CommonEditorComponent<T> extends RulesCommons implements CommonEditorContract<T>, OnInit {


    activatedRoute = inject(ActivatedRoute);
    notificationService = inject(NotificationService);
    service: any;
    router = inject(Router);
    formModel!: WritableSignal<T>;
    formGroup!: FieldTree<T>;

    constructor() {
        super();
        this.setService();
        this.setFormModel();
        this.setFormGroup();
    }

    ngOnInit(): void {
        this.getDetails();
    }

    setService(): void {
        throw new Error("Method not implemented.");
    }
    setFormModel(): void {
        throw new Error("Method not implemented.");
    }
    setFormGroup(): void {
        throw new Error("Method not implemented.");
    }

    getDetails() {
        this.setService();
        this.activatedRoute.queryParams.subscribe(params => {
            const editId = params['nimb_id'];
            const editVersion = params['version'];
            if (editId && editVersion) {
                this.service.get(editId, editVersion)
                    .subscribe((response: ApiResponse<T>) => {
                        this.formModel.set(response.data);
                        this.afterGetDetails();
                    });
            }
        });
    }

    saveDetails() {
        this.beforeSaveDetails();
        if (this.formGroup().valid()) {
            const formValue: T | any = this.formGroup().value();
            // Update existing variable package
            this.service.update(formValue.nimb_id, formValue.audit.version, formValue)
                .subscribe((response: ApiResponse<T>) => {
                    this.notificationService.success('Details updated successfully.', 5);
                });
        } else {
            this.notificationService.error('Please fill in all required fields.', 10);
        }
    }

    beforeSaveDetails(): void {
        // Optional hook for subclasses to implement additional logic before saving details
    }

    cancelEdit() {
        this.formGroup().reset();
        this.router.navigate([], {
            queryParams: {}
        });
    }

    afterGetDetails(): void {
        // Optional hook for subclasses to implement additional logic after getting details
    }

    navigateToVariablePackage(): void {
        const nimb_id = (<any>this.formModel()).variable_package.nimb_id;
        const version = (<any>this.formModel()).variable_package.audit.version;
        if (nimb_id && version) {
            const url = `/app/variable-packages?nimb_id=${encodeURIComponent(nimb_id)}&version=${encodeURIComponent(version)}`;
            window.open(url, '_blank');
        }
    }
}