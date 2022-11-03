import { Component, OnInit } from '@angular/core';

import { AppService } from '../app.service';
import { PendingApp } from '../pending-app';

@Component({
    selector: 'app-review',
    templateUrl: './review.component.html',
    styleUrls: ['./review.component.css']
})
export class ReviewComponent implements OnInit {
    apps: PendingApp[] = [];

    constructor(private appService: AppService) {}

    ngOnInit(): void {
        this.appService.getPendingApps().subscribe(apps => this.apps = apps);
    }

    approveApp(appId: string): void {
        this.appService.approveApp(appId).subscribe(_ => {
            const i = this.apps.findIndex(a => a.app_id === appId);
            if (i > -1) {
                this.apps.splice(i, 1);
            }
        });
    }
}
