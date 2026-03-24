#include "ELRDesktopAssistant.h"
#include <QApplication>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    ELRDesktopAssistant w;
    w.show();
    return a.exec();
}
