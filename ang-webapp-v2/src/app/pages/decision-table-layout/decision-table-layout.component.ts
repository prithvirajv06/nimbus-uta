import { Component } from '@angular/core';
import { CommonManagementComponent } from '../common-management-component';
import { DecisionTable } from '../../shared/types/dt.type';
import { DtTableComponent } from '../components/decision-table/dt-table/dt-table.component';
import { DtEditorComponent } from '../components/decision-table/dt-editor/dt-editor.component';
import { BasicTableComponent } from '../../shared/components/tables/basic-tables/basic-table.component';

@Component({
  selector: 'app-decision-table-layout',
  imports: [DtTableComponent, DtEditorComponent],
  templateUrl: './decision-table-layout.component.html',
  styleUrl: './decision-table-layout.component.css',
})
export class DecisionTableLayoutComponent extends CommonManagementComponent {
  decisionTable: DecisionTable | null = null;

  constructor() {
    super();
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

  edit(decisionTable: DecisionTable) {
    this.decisionTable = decisionTable;
    this.router.navigate([], {
      queryParams: {
        nimb_id: decisionTable.nimb_id,
        version: decisionTable.audit.version
      }
    });
  }


}
