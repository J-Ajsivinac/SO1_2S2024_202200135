#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h> // trae las funciones para crear archivos en /proc
#include <linux/seq_file.h> // trae las funciones para escribir en archivos en /proc
#include <linux/mm.h> // trae las funciones para manejar la memoria
#include <linux/sched.h> // trae las funciones para manejar los procesos
#include <linux/timer.h> // trae las funciones para manejar los timers
#include <linux/jiffies.h> // trae las funciones para manejar los jiffies, que son los ticks del sistema
#include <linux/mm.h>       // Para struct mm_struct
#include <linux/sched/stat.h> // Para el uso de CPU
#include <linux/sysinfo.h>
#include <linux/uaccess.h>  // Para copy_from_user

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tu Nombre");
MODULE_DESCRIPTION("Modulo para leer informacion de memoria y CPU");
MODULE_VERSION("1.0");

#define PROC_NAME "tarea1" // nombre del archivo en /proc

/* 
    Esta función se encarga de obtener la información de la memoria
    - si_meminfo: recibe un puntero a una estructura sysinfo, la cual se llena con la información de la memoria
*/
#include <linux/sched.h>    // Para for_each_process
#include <linux/mm.h>       // Para struct mm_struct
#include <linux/sched/stat.h> // Para el uso de CPU
#include <linux/seq_file.h>
#include <linux/sysinfo.h>
#include <linux/uaccess.h>  // Para copy_from_user


static int sysinfo_show(struct seq_file *m, void *v) {
    struct sysinfo si;
    struct task_struct *task;
    unsigned long rss, vsz;
    unsigned long total_ram_pages;

    si_meminfo(&si);

    total_ram_pages = totalram_pages();
    if (!total_ram_pages) {
        pr_err("No memory available\n");
        return -EINVAL;
    }

    unsigned long total_cpu_time = jiffies_to_msecs(get_jiffies_64());
    

    seq_printf(m, "{\n");
    seq_printf(m, "\t\"Total RAM\": %lu KB,\n", si.totalram / 1024);
    seq_printf(m, "\t\"Free RAM\": %lu KB,\n", si.freeram / 1024);
    seq_printf(m, "\t\"Used RAM\": %lu KB,\n", (si.totalram-si.freeram) / 1024);

    seq_printf(m, "\t\"Processes\": [\n");
    for_each_process(task) {
        if (strncmp(task->comm, "containerd", 10) == 0) {
            if(task->mm) {
                rss = get_mm_rss(task->mm) << (PAGE_SHIFT-10);
                vsz = task->mm->total_vm << (PAGE_SHIFT-10);
            } else {
                rss = 0;
                vsz = 0;
            }

            seq_printf(m, "\t\t\t\"PID\": %d,\n", task->pid);
            seq_printf(m, "\t\t\t\"Name\": \"%s\",\n", task->comm);
            // seq_printf(m, "\t\t\t\"VSZ\": %lu KB,\n", vsz);
            seq_printf(m, "\t\t\t\"RSS\": %lu KB,\n", rss);
            seq_printf(m, "\t\t\t\"Virtual Memory\": %lu KB,\n", vsz);
            unsigned long long percentage = (rss * 100ULL) / si.totalram;
            seq_printf(m, "\t\t\t\"Memory Usage\": %llu.%02llu%%,\n", percentage / 100, percentage % 100);
            unsigned long cpu_time = jiffies_to_msecs(task->utime + task->stime);
            unsigned long cpu_percentage = (cpu_time * 100) / total_cpu_time;
            seq_printf(m, "\t\t\t\"CPU\": %lu,\n", cpu_percentage);
        }
    }
        seq_printf(m, "\t]\n");
    seq_printf(m, "}\n");

    return 0;
}


// int read_proc(char *buf, char **start, off_t offset, int count, int *eof, void *data) {
//     int len = 0;
//     struct task_struct *task_list;
//     for_each_process(task_list) {
//         len += sprintf(buf+len, "PID: %d, Name: %s, State: %ld\n", task_list->pid, task_list->comm, task_list->state);
//     }
//     return len;
// }

/* 
    Esta función se ejecuta cuando se abre el archivo en /proc
    - single_open: se encarga de abrir el archivo y ejecutar la función sysinfo_show
*/
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

/* 
    Esta estructura contiene las operaciones a realizar cuando se accede al archivo en /proc
    - proc_open: se ejecuta cuando se abre el archivo
    - proc_read: se ejecuta cuando se lee el archivo
*/

static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
};


/* 
    Esta macro se encarga de hacer dos cosas:
    1. Ejecutar la función proc_create, la cual recibe el nombre del archivo a guardar en /proc, permisos,
        y la estructura con las operaciones a realizar

    2. Imprimir un mensaje en el log del kernel
*/
static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0, NULL, &sysinfo_ops);
    printk(KERN_INFO "tarea1 module loaded\n");
    return 0;
}

/* 
    Esta macro se encarga de hacer dos cosas:
    1. Ejecutar la función remove_proc_entry, la cual recibe el nombre del archivo a eliminar de /proc
*/
static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "tarea1 module unloaded\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);