import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewAppFormComponent } from './new-app-form.component';

describe('NewAppFormComponent', () => {
    let component: NewAppFormComponent;
    let fixture: ComponentFixture<NewAppFormComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [ NewAppFormComponent ]
        })
            .compileComponents();

        fixture = TestBed.createComponent(NewAppFormComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
