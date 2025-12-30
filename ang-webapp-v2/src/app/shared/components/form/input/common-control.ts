import { FormValueControl, ValidationError, WithOptionalField } from "@angular/forms/signals";
import { Component, Input, Output, EventEmitter, model, InputSignal, InputSignalWithTransform, ModelSignal, OutputRef, input } from '@angular/core';
@Component({
    selector: 'app-common-control',
    template: ``
})
export class CommonControlComponent implements FormValueControl<string | number | boolean> {

    value: ModelSignal<string | number | boolean> = model<string | number | boolean>('');

    disabled = input<boolean>(false);
    readonly = input<boolean>(false);
    hidden = input<boolean>(false);
    invalid = input<boolean>(false);
    touched = input<boolean>(false);
    errors = input<readonly WithOptionalField<ValidationError>[]>([]);
    dirty = input<boolean>(false);

    onChange(event: Event) {
        const val = (event.target as HTMLTextAreaElement).value;
        this.value.set(val);
    }
}