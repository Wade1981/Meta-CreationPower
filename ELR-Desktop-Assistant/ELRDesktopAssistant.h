#ifndef ELRDESKTOPASSISTANT_H
#define ELRDESKTOPASSISTANT_H

#include <QMainWindow>
#include "TrayIcon.h"
#include "DesktopWidget.h"
#include "ELRClient.h"

class ELRDesktopAssistant : public QMainWindow
{
    Q_OBJECT

public:
    ELRDesktopAssistant(QWidget *parent = nullptr);
    ~ELRDesktopAssistant();

private slots:
    void onELRStatusUpdated(const QString &status);
    void onELRContainersUpdated(const QList<QMap<QString, QString>> &containers);

private:
    TrayIcon *trayIcon;
    DesktopWidget *desktopWidget;
    ELRClient *elrClient;
    void initializeComponents();
    void setupConnections();
};

#endif // ELRDESKTOPASSISTANT_H
