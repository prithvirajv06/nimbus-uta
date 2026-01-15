import { Component, Directive, ElementRef, EventEmitter, HostListener, Injectable, Input, OnDestroy, Output, Renderer2 } from '@angular/core';
import { CommonControlComponent } from '../../form/input/common-control';
import { Option } from '../../form/select/select.component';
import { FormsModule } from '@angular/forms';
import { ClickOutsideDirective } from '../../../directives/clickoutside';

@Component({
  selector: 'app-search-select',
  standalone: true,
  imports: [FormsModule,ClickOutsideDirective],
  templateUrl: './search-select.component.html',
  styleUrls: ['./search-select.component.css'],
})
export class SearchSelectComponent extends CommonControlComponent {

  searchTerm: string = '';
  filterOptions: Option[] = [];
  @Input() options: Option[] | any = [];
  @Input() placeholder: string = 'Select an option';
  showDropdown: boolean = false;

  ngOnInit(): void {
    this.filterOptions = this.options;
  }

  onSearchTermChange(): void {
    const term = this.searchTerm.toLowerCase();
    this.filterOptions = this.options?.filter((option: Option) => option.label.toLowerCase().includes(term)) || [];
  }

  getLabelForValue(value: string | boolean | number): string {
    const foundOption = this.options?.find((option: Option) => option.value === value);
    return foundOption ? foundOption.label : '';
  }
}



