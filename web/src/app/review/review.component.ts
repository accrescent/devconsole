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
}
