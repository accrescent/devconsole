import { Component, OnInit } from '@angular/core';

import { App } from '../app';
import { AppService } from '../app.service';

@Component({
    selector: 'app-publish',
    templateUrl: './publish.component.html',
    styleUrls: ['./publish.component.css'],
})
export class PublishComponent implements OnInit {
    apps: App[] = [];

    constructor(private appService: AppService) {}

    ngOnInit(): void {
        this.appService.getApprovedApps().subscribe(apps => this.apps = apps);
    }

    publishApp(appId: string): void {
        this.appService.publishApp(appId).subscribe(_ => {
            const i = this.apps.findIndex(a => a.app_id === appId);
            if (i > -1) {
                this.apps.splice(i, 1);
            }
        });
    }
}
