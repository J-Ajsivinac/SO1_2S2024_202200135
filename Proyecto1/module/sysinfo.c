#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>q
#include <linux/sched.h>
#include <linux/mm.h>
#include <linux/uaccess.h>
#include <linux/slab.h>
#include <linux/cgroup.h>
#include <linux/fs.h>

#define FILE_NAME "sysinfo"
#define MAX_CMDLINE_LENGTH 1000


// Función para obtener la línea de comandos de un proceso
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

// Función para verificar si un proceso pertenece a un contenedor Docker
static int is_docker_container(struct task_struct *task) {
    // Verifica si el proceso padre es `containerd-shim`
    if (task && strstr(task->comm, "containerd-shim") != NULL) {
        return 1;
    }

    return 0;
}

// Función para capturar y mostrar solo el primer proceso del contenedor
// Función para capturar y mostrar información de todos los procesos del contenedor
static void get_container_processes_info(struct seq_file *m) {
    struct task_struct *task;
    bool found = false;

    for_each_process(task) {
        if (is_docker_container(task)) {
            struct mm_struct *mm = task->mm;
            unsigned long rss = 0, vsz = 0;

            if (mm) {
                rss = get_mm_rss(mm) * PAGE_SIZE / 1024;
                vsz = mm->total_vm * PAGE_SIZE / 1024;
            }

            seq_printf(m, "{\n\"pid\": %d,\n", task->pid);
            seq_printf(m, "\"name\": \"%s\",\n", get_process_cmdline(task));
            seq_printf(m, "\"cmdline\": \"%s\",\n", task->comm);
            seq_printf(m, "\"vsz\": %lu,\n", vsz);
            seq_printf(m, "\"rss\": %lu,\n", rss);
            seq_printf(m, "},\n");

            found = true;
        }
    }

    if (!found) {
        seq_printf(m, "{ \"error\": \"No container processes found\" }\n");
    }
}

// Función principal de secuencia para la lectura de /proc/sysinfo_#carnet
static int sysinfo_proc_show(struct seq_file *m, void *v) {
    seq_printf(m, "[\n");
    get_container_processes_info(m);
    seq_printf(m, "]\n");
    return 0;
}

// Funciones para abrir y mostrar el archivo en /proc
static int sysinfo_proc_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_proc_show, NULL);
}

static const struct proc_ops sysinfo_proc_ops = {
    .proc_open = sysinfo_proc_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

// Función de inicialización del módulo
static int __init sysinfo_module_init(void) {
    proc_create(FILE_NAME, 0, NULL, &sysinfo_proc_ops);
    return 0;
}

// Función de limpieza del módulo
static void __exit sysinfo_module_exit(void) {
    remove_proc_entry(FILE_NAME, NULL);
}

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tu Nombre");
MODULE_DESCRIPTION("Módulo de kernel para capturar información del primer proceso de un contenedor Docker en /proc");
MODULE_VERSION("1.0");

module_init(sysinfo_module_init);
module_exit(sysinfo_module_exit);
