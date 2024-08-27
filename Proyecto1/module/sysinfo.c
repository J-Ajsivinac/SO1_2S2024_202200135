#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched.h>
#include <linux/mm.h>
#include <linux/uaccess.h>
#include <linux/slab.h>
#include <linux/cgroup.h>
#include <linux/fs.h>
// sysinfo_202200135
#define FILE_NAME "sysinfo"
#define MAX_CMDLINE_LENGTH 1000


static char *get_process_cmdline(struct task_struct *task) {
    struct mm_struct *mm;
    char *cmdline, *p;
    unsigned long arg_start, arg_end, env_start;
    int i, len;

    cmdline = kmalloc(MAX_CMDLINE_LENGTH, GFP_KERNEL);
    if (!cmdline)
        return NULL;

    mm = get_task_mm(task);
    if (!mm) {
        kfree(cmdline);
        return NULL;
    }

    down_read(&mm->mmap_lock);
    arg_start = mm->arg_start;
    arg_end = mm->arg_end;
    env_start = mm->env_start;
    up_read(&mm->mmap_lock);

    len = arg_end - arg_start;

    if (len > MAX_CMDLINE_LENGTH - 1)
        len = MAX_CMDLINE_LENGTH - 1;

    if (access_process_vm(task, arg_start, cmdline, len, 0) != len) {
        mmput(mm);
        kfree(cmdline);
        return NULL;
    }

    cmdline[len] = '\0';

    // Reemplazar caracteres nulos por espacios
    p = cmdline;
    for (i = 0; i < len; i++)
        if (p[i] == '\0')
            p[i] = ' ';

    mmput(mm);
    return cmdline;
}

unsigned total = 0, used = 0, free_r = 0;

static void get_memory_info(struct seq_file *m){
    struct sysinfo i;
    si_meminfo(&i);

    unsigned long toal_ram = i.totalram * i.mem_unit;
    total = toal_ram;
    unsigned long free_ram = i.freeram * i.mem_unit;
    free_r = free_ram;
    unsigned long used_ram = toal_ram - free_ram;
    used = used_ram;
    seq_printf(m, "Memory:\n");
    seq_printf(m, "{\n\"total_ram\": %lu,\n", toal_ram / 1024);
    seq_printf(m, "\"free_ram\": %lu,\n", free_ram / 1024);
    seq_printf(m, "\"used_ram\": %lu,\n", used_ram / 1024);
    seq_printf(m, "},\n");
}

//Función para verificar si un proceso pertenece a un contenedor Docker
static int is_docker_container(struct task_struct *task) {
    // Verifica si el proceso padre es `containerd-shim`
    if (task && strstr(task->comm, "containerd-shim") != NULL) {
        return 1;
    }

    return 0;
}


static void get_container_processes_info(struct seq_file *m) {
    struct task_struct *task;
    bool found = false;

    struct sysinfo i;
    si_meminfo(&i);
    signed long toal_ram = i.totalram * i.mem_unit;
    unsigned long total_jiffies = jiffies;
    unsigned long cpu_usage;
    unsigned long cpu_usage_child;

    // unsigned long total_jiffies = jiffies;
    for_each_process(task) {
        if (is_docker_container(task)==1) {
            // struct sched_entity *se = &task->se;
            struct mm_struct *mm = task->mm;
            unsigned long rss = 0, vsz = 0;
            unsigned long rssKB = 0, vszKB = 0;
            unsigned long porc_ram = 0;
            if (mm) {
                rss = get_mm_rss(mm) << PAGE_SHIFT;
                // rssKB = rss / 1024;
                vsz = mm->total_vm << PAGE_SHIFT;
                // vszKB = vsz / 1024;
            }
            unsigned long total_ram_pages;
            total_ram_pages = totalram_pages();
            if(found){
                seq_printf(m, ",\n");
            }

            unsigned long total_time = task->utime + task->stime;
            cpu_usage = (total_time * 10000) / total_jiffies;
                
            
            if(task->children.next != NULL){
                struct task_struct *child;
                list_for_each_entry(child, &task->children, sibling){
                    struct mm_struct *mm_child = child->mm;
                    if(mm_child){
                        rss += get_mm_rss(mm_child) << PAGE_SHIFT;
                        vsz += mm_child->total_vm << PAGE_SHIFT;
                    }
                    unsigned long total_time_child = child->utime + child->stime;
                    // total_time_child /= HZ;
                    cpu_usage_child = (total_time_child * 10) / (total_jiffies);
                    cpu_usage += cpu_usage_child;
                    // cpu_usage = (total_time_child * 10000) / (total_jiffies);
                    // seq_printf(m, "%s", get_process_cmdline(child));
                }
                
            }
            rssKB = rss / 1024;
            vszKB = vsz / 1024;

            if(found){
                seq_printf(m, ",\n");
            }
            
            seq_printf(m, "{\n");
            seq_printf(m, "\"pid\": %d,\n", task->pid);
            seq_printf(m, "\"name\": \"%s\",\n", get_process_cmdline(task));
            seq_printf(m, "\"cmdline\": \"%s\",\n", task->comm);
            seq_printf(m, "\"vsz\": %lu,\n", vszKB);
            seq_printf(m, "\"rss\": %lu,\n", rssKB);
            porc_ram = (rss * 10000) / toal_ram;
            seq_printf(m, "\"mem percent\": %lu.%02lu,\n", porc_ram/100, porc_ram%100);
            seq_printf(m, "\"CPUUsage\": %lu.%02lu\n", cpu_usage / 100, cpu_usage % 100);
            // seq_printf(m, "\"CPUUsageChild\": %lu.%02lu\n", cpu_usage_child / 100, cpu_usage_child % 100);
            seq_printf(m, "}");
            found = true;
        }
    }

    if (!found) {
        seq_printf(m, "{ \"error\": \"No container processes found\" }\n");
    }
}

static int sysinfo_proc_show(struct seq_file *m, void *v) {
    seq_printf(m, "{\n");
    get_memory_info(m);
    seq_printf(m, "Processes:\n");
    seq_printf(m, "[\n");
    get_container_processes_info(m);
    seq_printf(m, "]\n");
    return 0;
}

static int sysinfo_proc_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_proc_show, NULL);
}

static const struct proc_ops sysinfo_proc_ops = {
    .proc_open = sysinfo_proc_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

static int __init sysinfo_module_init(void) {
    proc_create(FILE_NAME, 0, NULL, &sysinfo_proc_ops);
    return 0;
}


static void __exit sysinfo_module_exit(void) {
    remove_proc_entry(FILE_NAME, NULL);
}

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Joab Ajsivianc");
MODULE_DESCRIPTION("Módulo de kernel para capturar información de los procesos de un contenedor Docker en /proc");
MODULE_VERSION("1.0");

module_init(sysinfo_module_init);
module_exit(sysinfo_module_exit);
