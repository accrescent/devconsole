import { TestBed } from '@angular/core/testing';

import { ReviewerGuard } from './reviewer.guard';

describe('ReviewerGuard', () => {
    let guard: ReviewerGuard;

    beforeEach(() => {
        TestBed.configureTestingModule({});
        guard = TestBed.inject(ReviewerGuard);
    });

    it('should be created', () => {
        expect(guard).toBeTruthy();
    });
});
