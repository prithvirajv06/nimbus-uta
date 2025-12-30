import { Component } from '@angular/core';
import { CommonManagementComponent } from '../common-management-component';
import { VariablePackage } from '../../shared/types/variable_package';
import { VariablePackageEditorComponent } from '../components/variable-package/variable-package-editor/variable-package-editor.component';
import { VariablePackageTableComponent } from '../components/variable-package/variable-package-table/variable-package-table.component';

@Component({
  selector: 'app-variable-package',
  imports: [VariablePackageTableComponent, VariablePackageEditorComponent],
  templateUrl: './variable-package-layout.component.html',
  styleUrl: './variable-package-layout.component.css',
})
export class VariablePackageLayoutComponent extends CommonManagementComponent {
  variablePack: VariablePackage | null = null;

  constructor() {
    super();
    this.activatedRoute.queryParams.subscribe(params => {
      const editId = params['nimb_id'];
      const editVersion = params['version'];
      if (editId && editVersion) {
        this.viewMode.set('editor');
      }else{
        this.viewMode.set('list');
      }
    });
  }
  
  editVariablePackage(variablePack: VariablePackage) {
    this.variablePack = variablePack;
    this.router.navigate([], {
      queryParams: {
        nimb_id: variablePack.nimb_id,
        version: variablePack.audit.version
      }
    });
  }


}
