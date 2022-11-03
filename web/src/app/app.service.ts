import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs';

import { App } from './app';
import { PendingApp } from './pending-app';

@Injectable({
    providedIn: 'root'
})
export class AppService {
    private readonly appsUrl = 'api/apps';
    private readonly pendingAppsUrl = 'api/pending-apps';

    constructor(private http: HttpClient) {}

    getPendingApps(): Observable<PendingApp[]> {
        return this.http.get<PendingApp[]>(this.pendingAppsUrl);
    }

    approveApp(appId: string): Observable<void> {
        return this.http.patch<void>(`${this.pendingAppsUrl}/${appId}`, '');
    }

    rejectApp(appId: string): Observable<void> {
        return this.http.delete<void>(`${this.pendingAppsUrl}/${appId}`);
    }

    uploadApp(app: File): Observable<App> {
        const formData = new FormData();
        formData.append("app", app);

        return this.http.post<App>(this.appsUrl, formData);
    }

    submitApp(id: string): Observable<void> {
        return this.http.patch<void>(`${this.appsUrl}/${id}`, '');
    }
}
