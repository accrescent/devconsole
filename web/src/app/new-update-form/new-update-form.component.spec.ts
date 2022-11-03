import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewUpdateFormComponent } from './new-update-form.component';

describe('NewUpdateFormComponent', () => {
    let component: NewUpdateFormComponent;
    let fixture: ComponentFixture<NewUpdateFormComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [ NewUpdateFormComponent ]
        })
            .compileComponents();

        fixture = TestBed.createComponent(NewUpdateFormComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
