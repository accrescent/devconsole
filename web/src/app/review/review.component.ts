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
    updates: PendingApp[] = [];

    constructor(private appService: AppService) {}

    ngOnInit(): void {
        this.appService.getPendingApps().subscribe(apps => this.apps = apps);
        this.appService.getUpdates().subscribe(updates => this.updates = updates);
    }

    approveApp(appId: string): void {
        this.appService.approveApp(appId).subscribe(_ => this.removeApp(appId));
    }

    rejectApp(appId: string): void {
        this.appService.rejectApp(appId).subscribe(_ => this.removeApp(appId));
    }

    approveUpdate(appId: string, versionCode: number): void {
        this.appService.approveUpdate(appId, versionCode)
            .subscribe(_ => this.removeUpdate(appId, versionCode));
    }

    rejectUpdate(appId: string, versionCode: number): void {
        this.appService.rejectUpdate(appId, versionCode)
            .subscribe(_ => this.removeUpdate(appId, versionCode));
    }


    private removeApp(appId: string): void {
        const i = this.apps.findIndex(a => a.app_id === appId);
        if (i > -1) {
            this.apps.splice(i, 1);
        }
    }

    private removeUpdate(appId: string, versionCode: number): void {
        const i = this.updates.findIndex(u => u.app_id === appId && u.version_code === versionCode);
        if (i > -1) {
            this.updates.splice(i, 1);
        }
    }
}
