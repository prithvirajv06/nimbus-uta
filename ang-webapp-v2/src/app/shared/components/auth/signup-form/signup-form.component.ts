
import { Component, signal, inject } from '@angular/core';
import { email, form, required, Field } from '@angular/forms/signals';
import { LabelComponent } from '../../form/label/label.component';
import { CheckboxComponent } from '../../form/input/checkbox.component';
import { InputFieldComponent } from '../../form/input/input-field.component';
import { Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AuthenticationService } from '../../../services/authentication.service';
import { AppCustomer } from '../../../types/user.type';
import { ApiResponse } from '../../../types/common.type';
import { NotificationService } from '../../../services/notification.service';


@Component({
  selector: 'app-signup-form',
  imports: [
    LabelComponent,
    CheckboxComponent,
    InputFieldComponent,
    RouterModule,
    FormsModule,
    Field
  ],
  templateUrl: './signup-form.component.html',
  styles: ``
})
export class SignupFormComponent {

  authService = inject(AuthenticationService)
  notificationService = inject(NotificationService)
  router = inject(Router)

  showPassword = false;
  isChecked = false;

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
    required(schemaPath.fname, { message: 'First Name is required' });
    required(schemaPath.lname, { message: 'Last Name is required' });
    required(schemaPath.organization.name, { message: 'Organisation is required' });
    required(schemaPath.email, { message: 'Email is required' });
    email(schemaPath.email, { message: 'Enter a valid email address' });
    required(schemaPath.password, { message: 'Password is required' });
  })


  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  onSignIn() {
    console.log(this.formGroup().value());
    if(this.formGroup().invalid()){
      this.notificationService.error('Please fill all required fields correctly.', 5);
      return;
    }
    if(!this.isChecked){
      this.notificationService.error('You must agree to the Terms and Conditions.', 5);
      return;
    }
    this.authService.signUpUser(<AppCustomer>this.formGroup().value()).subscribe((response:ApiResponse<AppCustomer>) => {
      this.notificationService.success(response.message, 5);
      this.formGroup().reset();
      this.router.navigate(['/auth/sign-in']);
    });
  }
}
