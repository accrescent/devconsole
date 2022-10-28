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
    form = this.fb.group({});

    constructor(
        private fb: NonNullableFormBuilder,
        private appService: AppService,
        private router: Router,
    ) {}

    onFileChange(event: Event): void {
        const file = (event.target as HTMLInputElement).files?.[0];

        if (file !== undefined) {
            this.appService.uploadApp(file).subscribe(app => this.app = app);
        }
    }

    onSubmit(): void {
        this.appService.submitApp(this.app!.id).subscribe(_ => this.router.navigate(['console']));
    }
}
