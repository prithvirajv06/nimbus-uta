import { NgClass } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-side-draw',
  imports: [NgClass],
  templateUrl: './side-draw.component.html',
  styleUrl: './side-draw.component.css',
})
export class SideDrawComponent {

  @Input() title: string = 'Drawer Menu';
  @Input() className: string = '';
  @Output() close = new EventEmitter<void>();
  @Input() isOpen: boolean = false;
}
