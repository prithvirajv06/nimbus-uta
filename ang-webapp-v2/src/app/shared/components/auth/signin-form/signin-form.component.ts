
import { Component, inject, signal } from '@angular/core';
import { LabelComponent } from '../../form/label/label.component';
import { CheckboxComponent } from '../../form/input/checkbox.component';
import { ButtonComponent } from '../../ui/button/button.component';
import { InputFieldComponent } from '../../form/input/input-field.component';
import { Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AuthenticationService } from '../../../services/authentication.service';
import { NotificationService } from '../../../services/notification.service';
import { form, required, email, Field } from '@angular/forms/signals';
import { AppCustomer } from '../../../types/user.type';
import { ApiResponse } from '../../../types/common.type';

@Component({
  selector: 'app-signin-form',
  imports: [
    LabelComponent,
    CheckboxComponent,
    ButtonComponent,
    InputFieldComponent,
    RouterModule,
    FormsModule,
    Field
  ],
  templateUrl: './signin-form.component.html',
  styles: ``
})
export class SigninFormComponent {

  showPassword = false;
  isChecked = false;
  authService = inject(AuthenticationService)
  notificationService = inject(NotificationService)
  router = inject(Router)

  formModel = signal<AppCustomer>({
    fname: '',
    lname: '',
    organization: {
      nimb_id: '',
      name: '',
      address: ''
    },
    email: '',
    password: '',
    nimb_id: '',
  })
  formGroup = form(this.formModel, (schemaPath) => {
    required(schemaPath.email, { message: 'Email is required' });
    email(schemaPath.email, { message: 'Enter a valid email address' });
    required(schemaPath.password, { message: 'Password is required' });
  })

  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  onSignIn() {
    this.authService.singInUser(this.formModel().email, this.formModel().password).subscribe((response: ApiResponse<AppCustomer>) => {
      this.notificationService.success(response.message, 5);
      this.formGroup().reset();
      this.authService.clearCurrentUser();
      this.authService.setCurrentUser(response.data!);
      this.router.navigate(['/app/dashboard']);
    });
  }
}
