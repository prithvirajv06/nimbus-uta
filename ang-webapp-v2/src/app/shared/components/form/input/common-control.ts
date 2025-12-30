import { FormValueControl, ValidationError, WithOptionalField } from "@angular/forms/signals";
import { Component, Input, Output, EventEmitter, model, InputSignal, InputSignalWithTransform, ModelSignal, OutputRef, input } from '@angular/core';
@Component({
    selector: 'app-common-control',
    template: ``
})
export class CommonControlComponent implements FormValueControl<string | number> {

    value: ModelSignal<string | number> = model<string | number>('');

    disabled = input<boolean>(false);
    readonly = input<boolean>(false);
    hidden = input<boolean>(false);
    invalid = input<boolean>(false);
    touched = input<boolean>(false);
    errors = input<readonly WithOptionalField<ValidationError>[]>([]);
    dirty = input<boolean>(false);

    onInput(event: Event) {
        const val = (event.target as HTMLTextAreaElement).value;
        this.value.set(val);
    }
}