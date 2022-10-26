import { Component, OnInit } from '@angular/core';
import { NonNullableFormBuilder, Validators } from '@angular/forms';

import { RegisterService } from '../register.service';
import { Email } from '../email';

@Component({
    selector: 'app-register-form',
    templateUrl: './register-form.component.html',
    styleUrls: ['./register-form.component.css']
})
export class RegisterFormComponent implements OnInit {
    form = this.fb.group({
        email: this.fb.control('', [Validators.required])
    });
    emails: string[] = [];

    constructor(private fb: NonNullableFormBuilder, private registerService: RegisterService) {}

    ngOnInit(): void {
        this.registerService.getEmails().subscribe(emails =>
            emails.forEach((email, _) => this.emails.push(email))
        )
    }

    onSubmit(): void {
        const email: string = this.form.getRawValue().email;
        this.registerService.register({ email }).subscribe();
    }
}
