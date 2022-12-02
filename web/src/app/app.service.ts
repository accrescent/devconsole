import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs';

import { App } from './app';

@Injectable({
    providedIn: 'root'
})
export class AppService {
    private readonly appsUrl = 'api/apps';
    private readonly submittedAppsUrl = 'api/submitted-apps';
    private readonly submittedUpdatesUrl = 'api/submitted-updates';
    private readonly pendingAppsUrl = 'api/pending-apps';
    private readonly updatesUrl = 'api/updates';
    private readonly approvedAppsUrl = 'api/approved-apps';

    constructor(private http: HttpClient) {}

    getApps(): Observable<App[]> {
        return this.http.get<App[]>(this.appsUrl);
    }

    getSubmittedApps(): Observable<App[]> {
        return this.http.get<App[]>(this.submittedAppsUrl);
    }

    getSubmittedUpdates(): Observable<App[]> {
        return this.http.get<App[]>(this.submittedUpdatesUrl);
    }

    getPendingApps(): Observable<App[]> {
        return this.http.get<App[]>(this.pendingAppsUrl);
    }

    getUpdates(): Observable<App[]> {
        return this.http.get<App[]>(this.updatesUrl);
    }

    getApprovedApps(): Observable<App[]> {
        return this.http.get<App[]>(this.approvedAppsUrl);
    }

    approveApp(appId: string): Observable<void> {
        return this.http.patch<void>(`${this.pendingAppsUrl}/${appId}`, '');
    }

    rejectApp(appId: string): Observable<void> {
        return this.http.delete<void>(`${this.pendingAppsUrl}/${appId}`);
    }

    uploadApp(app: File, icon: File): Observable<App> {
        const formData = new FormData();
        formData.append("app", app);
        formData.append("icon", icon);

        return this.http.post<App>(this.appsUrl, formData);
    }

    uploadUpdate(app: File, appId: string): Observable<App> {
        const formData = new FormData();
        formData.append("app", app);

        return this.http.post<App>(`${this.appsUrl}/${appId}/updates`, formData);
    }

    submitApp(id: string, label: string): Observable<void> {
        return this.http.patch<void>(`${this.appsUrl}/${id}`, { label });
    }

    submitUpdate(id: string, versionCode: number): Observable<void> {
        return this.http.patch<void>(`${this.appsUrl}/${id}/${versionCode}`, '');
    }

    approveUpdate(id: string, versionCode: number): Observable<void> {
        return this.http.patch<void>(`${this.updatesUrl}/${id}/${versionCode}`, '');
    }

    rejectUpdate(id: string, versionCode: number): Observable<void> {
        return this.http.delete<void>(`${this.updatesUrl}/${id}/${versionCode}`);
    }

    publishApp(id: string): Observable<void> {
        return this.http.post<void>(`${this.appsUrl}/${id}`, '');
    }
}
