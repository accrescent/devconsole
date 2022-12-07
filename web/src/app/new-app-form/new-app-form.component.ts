import { Component } from '@angular/core';
import { HttpEventType, HttpResponse } from '@angular/common/http';
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
    uploadProgress = 0;
    confirmationForm = this.fb.group({
        label: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(30)]],
    });

    constructor(
        private fb: NonNullableFormBuilder,
        private appService: AppService,
        private router: Router,
    ) {}

    onUpload(): void {
        const app = (<HTMLInputElement>document.getElementById("app")).files?.[0];
        const icon = (<HTMLInputElement>document.getElementById("icon")).files?.[0];

        if (app !== undefined && icon !== undefined) {
            this.appService.uploadApp(app, icon).subscribe(event => {
                if (event.type === HttpEventType.UploadProgress) {
                    this.uploadProgress = 100 * event.loaded / event.total!!;

                    // Clear the progress bar once the upload is complete
                    if (event.loaded === event.total!!) {
                        this.uploadProgress = 0;
                    }
                } else if (event instanceof HttpResponse) {
                    this.app = event.body!!;
                    this.confirmationForm.patchValue({ label: this.app.label });
                }
            });
        }
    }

    onConfirm(): void {
        const label = this.confirmationForm.getRawValue().label;
        this.appService.submitApp(this.app!.app_id, label)
            .subscribe(_ => this.router.navigate(['dashboard']));
    }
}
