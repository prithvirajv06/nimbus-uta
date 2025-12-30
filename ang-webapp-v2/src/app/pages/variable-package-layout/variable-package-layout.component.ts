import { Component } from '@angular/core';
import { CommonLayoutComponent } from '../common-management-component';
import { VariablePackage } from '../../shared/types/variable_package';
import { VariablePackageEditorComponent } from '../components/variable-package/variable-package-editor/variable-package-editor.component';
import { VariablePackageTableComponent } from '../components/variable-package/variable-package-table/variable-package-table.component';

@Component({
  selector: 'app-variable-package',
  imports: [VariablePackageTableComponent, VariablePackageEditorComponent],
  templateUrl: './variable-package-layout.component.html',
  styleUrl: './variable-package-layout.component.css',
})
export class VariablePackageLayoutComponent extends CommonLayoutComponent<VariablePackage> {



}
