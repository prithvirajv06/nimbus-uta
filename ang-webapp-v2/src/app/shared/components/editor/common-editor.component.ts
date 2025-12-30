import { inject, WritableSignal } from "@angular/core";
import { FieldTree } from "@angular/forms/signals";
import { ActivatedRoute, Router } from "@angular/router";
import { NotificationService } from "../../services/notification.service";
import { ApiResponse } from "../../types/common.type";

export interface CommonEditorContract<T> {
    setService(): void;
    setFormModel(): void;
    setFormGroup(): void;
    saveDetails(): void;
    getDetails(): void;
}

export class CommonEditorComponent<T> implements CommonEditorContract<T> {


    activatedRoute = inject(ActivatedRoute);
    notificationService = inject(NotificationService);
    service: any;
    router = inject(Router);
    formModel!: WritableSignal<T>;
    formGroup!: FieldTree<T>;

    constructor() {

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
        this.activatedRoute.queryParams.subscribe(params => {
          const editId = params['nimb_id'];
          const editVersion = params['version'];
          if (editId && editVersion) {
            this.service.get(editId, editVersion)
              .subscribe((response: ApiResponse<T>) => {
                this.formModel.set(response.data);
              });
          }
        });
      }

    saveDetails() {
        this.service.update(this.formGroup().value()).subscribe({
            next: (res: ApiResponse<T>) => {
                this.notificationService.success('Details saved successfully');
            }
        });
    }
}