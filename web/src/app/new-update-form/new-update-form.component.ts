import { Component, OnInit } from '@angular/core';
import { HttpEventType, HttpResponse } from '@angular/common/http';
import { NonNullableFormBuilder } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { App } from '../app';
import { AppService } from '../app.service';

@Component({
    selector: 'app-new-update-form',
    templateUrl: './new-update-form.component.html',
    styleUrls: ['./new-update-form.component.css']
})
export class NewUpdateFormComponent implements OnInit {
    private appId = "";
    app: App | undefined = undefined;
    form = this.fb.group({});
    uploadProgress = 0;

    constructor(
        private fb: NonNullableFormBuilder,
        private appService: AppService,
        private router: Router,
        private activatedRoute: ActivatedRoute,
    ) {}

    ngOnInit(): void {
        this.activatedRoute.paramMap.subscribe(params => this.appId = params.get('id') ?? "");
    }

    onFileChange(event: Event): void {
        const file = (event.target as HTMLInputElement).files?.[0];

        if (file !== undefined) {
            this.appService.uploadUpdate(file, this.appId).subscribe(event => {
                if (event.type === HttpEventType.UploadProgress) {
                    this.uploadProgress = 100 * event.loaded / event.total!!;

                    // Clear the progress bar once the upload is complete
                    if (event.loaded === event.total!!) {
                        this.uploadProgress = 0;
                    }
                } else if (event instanceof HttpResponse) {
                    this.app = event.body!!;
                }
            });
        }
    }

    onSubmit(): void {
        this.appService.submitUpdate(this.app!.app_id, this.app!.version_code)
            .subscribe(_ => this.router.navigate(['apps']));
    }
}
