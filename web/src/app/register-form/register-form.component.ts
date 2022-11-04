import { Component, OnInit } from '@angular/core';
import { NonNullableFormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { RegisterService } from '../register.service';

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

    constructor(
        private fb: NonNullableFormBuilder,
        private registerService: RegisterService,
        private router: Router,
    ) {}

    ngOnInit(): void {
        this.registerService.getEmails().subscribe(emails =>
            emails.forEach((email, _) => this.emails.push(email))
        )
    }

    onSubmit(): void {
        const email: string = this.form.getRawValue().email;
        this.registerService.register({ email }).subscribe(_ =>
            this.router.navigate(['dashboard'])
        );
    }
}
