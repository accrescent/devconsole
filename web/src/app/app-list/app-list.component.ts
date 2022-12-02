import { Component, OnInit } from '@angular/core';

import { App } from '../app';
import { AppService } from '../app.service';

@Component({
    selector: 'app-app-list',
    templateUrl: './app-list.component.html',
    styleUrls: ['./app-list.component.css']
})
export class AppListComponent implements OnInit {
    apps: App[] = [];
    submittedApps: App[] = [];
    submittedUpdates: App[] = [];

    constructor(private appService: AppService) {}

    ngOnInit(): void {
        this.appService.getApps().subscribe(apps => this.apps = apps);
        this.appService.getSubmittedApps().subscribe(apps => this.submittedApps = apps);
        this.appService.getSubmittedUpdates().subscribe(updates => this.submittedUpdates = updates);
    }
}
