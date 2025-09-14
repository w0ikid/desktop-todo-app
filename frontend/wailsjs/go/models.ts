export namespace app {
	
	export class CreateTaskOutput {
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateTaskOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class GetDashboardOutput {
	    active_count: number;
	    completed_count: number;
	    overdue_count: number;
	    due_today: domain.Task[];
	    due_this_week: domain.Task[];
	    recent_tasks: domain.Task[];
	
	    static createFrom(source: any = {}) {
	        return new GetDashboardOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.active_count = source["active_count"];
	        this.completed_count = source["completed_count"];
	        this.overdue_count = source["overdue_count"];
	        this.due_today = this.convertValues(source["due_today"], domain.Task);
	        this.due_this_week = this.convertValues(source["due_this_week"], domain.Task);
	        this.recent_tasks = this.convertValues(source["recent_tasks"], domain.Task);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetTaskOutput {
	    task?: domain.Task;
	
	    static createFrom(source: any = {}) {
	        return new GetTaskOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task = this.convertValues(source["task"], domain.Task);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListTasksOutput {
	    tasks: domain.Task[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ListTasksOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tasks = this.convertValues(source["tasks"], domain.Task);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace domain {
	
	export class Task {
	    ID: string;
	    Title: string;
	    Status: string;
	    CreatedAt: time.Time;
	    DueDate?: time.Time;
	    Priority: string;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Title = source["Title"];
	        this.Status = source["Status"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], time.Time);
	        this.DueDate = this.convertValues(source["DueDate"], time.Time);
	        this.Priority = source["Priority"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace time {
	
	export class Time {
	
	
	    static createFrom(source: any = {}) {
	        return new Time(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

