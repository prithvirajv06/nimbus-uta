import { Component } from '@angular/core';
import { CommonLayoutComponent } from '../common-management-component';
import { DecisionTable } from '../../shared/types/dt.type';
import { DtTableComponent } from '../components/decision-table/dt-table/dt-table.component';
import { DtEditorComponent } from '../components/decision-table/dt-editor/dt-editor.component';

@Component({
  selector: 'app-decision-table-layout',
  imports: [DtTableComponent, DtEditorComponent],
  templateUrl: './decision-table-layout.component.html',
  styleUrl: './decision-table-layout.component.css',
})
export class DecisionTableLayoutComponent extends CommonLayoutComponent<DecisionTable> {


}
