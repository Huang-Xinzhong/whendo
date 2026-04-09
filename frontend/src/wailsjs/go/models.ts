export namespace models {
	
	export class Task {
	    id: number;
	    workspace_id: number;
	    title: string;
	    description: string;
	    type: string;
	    // Go type: time
	    due_at?: any;
	    // Go type: time
	    remind_at?: any;
	    is_completed: boolean;
	    start_time?: string;
	    end_time?: string;
	    interval_value?: number;
	    interval_unit?: string;
	    repeat_mode?: string;
	    weekdays?: string;
	    month_day?: number;
	    // Go type: time
	    next_trigger_at?: any;
	    // Go type: time
	    paused_date?: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    remindText?: string;
	    pausedToday?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workspace_id = source["workspace_id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.due_at = this.convertValues(source["due_at"], null);
	        this.remind_at = this.convertValues(source["remind_at"], null);
	        this.is_completed = source["is_completed"];
	        this.start_time = source["start_time"];
	        this.end_time = source["end_time"];
	        this.interval_value = source["interval_value"];
	        this.interval_unit = source["interval_unit"];
	        this.repeat_mode = source["repeat_mode"];
	        this.weekdays = source["weekdays"];
	        this.month_day = source["month_day"];
	        this.next_trigger_at = this.convertValues(source["next_trigger_at"], null);
	        this.paused_date = this.convertValues(source["paused_date"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.remindText = source["remindText"];
	        this.pausedToday = source["pausedToday"];
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
	export class Workspace {
	    id: number;
	    name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    color?: string;
	    taskCount?: number;
	
	    static createFrom(source: any = {}) {
	        return new Workspace(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.color = source["color"];
	        this.taskCount = source["taskCount"];
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

export namespace requests {
	
	export class SettingsUpdateReq {
	    default_workspace_id?: string;
	    default_sort?: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingsUpdateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.default_workspace_id = source["default_workspace_id"];
	        this.default_sort = source["default_sort"];
	    }
	}
	export class TaskCreateReq {
	    workspace_id: number;
	    title: string;
	    description: string;
	    type: string;
	    due_at?: string;
	    remind_at?: string;
	    start_time?: string;
	    end_time?: string;
	    interval_value?: number;
	    interval_unit?: string;
	    repeat_mode?: string;
	    weekdays?: string;
	    month_day?: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskCreateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.workspace_id = source["workspace_id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.due_at = source["due_at"];
	        this.remind_at = source["remind_at"];
	        this.start_time = source["start_time"];
	        this.end_time = source["end_time"];
	        this.interval_value = source["interval_value"];
	        this.interval_unit = source["interval_unit"];
	        this.repeat_mode = source["repeat_mode"];
	        this.weekdays = source["weekdays"];
	        this.month_day = source["month_day"];
	    }
	}
	export class TaskUpdateReq {
	    id: number;
	    workspace_id?: number;
	    title?: string;
	    description?: string;
	    type?: string;
	    due_at?: string;
	    remind_at?: string;
	    start_time?: string;
	    end_time?: string;
	    interval_value?: number;
	    interval_unit?: string;
	    repeat_mode?: string;
	    weekdays?: string;
	    month_day?: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskUpdateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workspace_id = source["workspace_id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.due_at = source["due_at"];
	        this.remind_at = source["remind_at"];
	        this.start_time = source["start_time"];
	        this.end_time = source["end_time"];
	        this.interval_value = source["interval_value"];
	        this.interval_unit = source["interval_unit"];
	        this.repeat_mode = source["repeat_mode"];
	        this.weekdays = source["weekdays"];
	        this.month_day = source["month_day"];
	    }
	}
	export class WorkspaceCreateReq {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkspaceCreateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class WorkspaceUpdateReq {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkspaceUpdateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

