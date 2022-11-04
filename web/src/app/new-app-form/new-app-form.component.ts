import { Component } from '@angular/core';
import { NonNullableFormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { App } from '../app';
import { AppService } from '../app.service';

@Component({
    selector: 'app-new-app-form',
    templateUrl: './new-app-form.component.html',
    styleUrls: ['./new-app-form.component.css']
})
export class NewAppFormComponent {
    app: App | undefined = undefined;
    uploadForm = this.fb.group({
        app: ['', Validators.required],
        icon: ['', Validators.required],
    });
    confirmationForm = this.fb.group({});

    constructor(
        private fb: NonNullableFormBuilder,
        private appService: AppService,
        private router: Router,
    ) {}

    onUpload(): void {
        const app = (<HTMLInputElement>document.getElementById("app")).files?.[0];
        const icon = (<HTMLInputElement>document.getElementById("icon")).files?.[0];

        if (app !== undefined && icon !== undefined) {
            this.appService.uploadApp(app, icon).subscribe(app => this.app = app);
        }
    }

    onConfirm(): void {
        this.appService.submitApp(this.app!.app_id)
            .subscribe(_ => this.router.navigate(['dashboard']));
    }
}
