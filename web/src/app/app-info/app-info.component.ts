import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
    selector: 'app-app-info',
    templateUrl: './app-info.component.html',
    styleUrls: ['./app-info.component.css'],
})
export class AppInfoComponent implements OnInit {
    appId = "";

    constructor(private activatedRoute: ActivatedRoute) {}

    ngOnInit(): void {
        this.activatedRoute.paramMap.subscribe(params => this.appId = params.get('id') ?? "");
    }
}
