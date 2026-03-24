#ifndef TRAYICON_H
#define TRAYICON_H

#include <QObject>
#include <QSystemTrayIcon>
#include <QMenu>

class TrayIcon : public QObject
{
    Q_OBJECT

public:
    explicit TrayIcon(QObject *parent = nullptr);
    ~TrayIcon();

public slots:
    void updateStatus(const QString &status);

signals:
    void showDesktopWidget();
    void hideDesktopWidget();
    void exitApp();

private slots:
    void onTrayIconActivated(QSystemTrayIcon::ActivationReason reason);
    void onShowWidget();
    void onHideWidget();
    void onExit();

private:
    QSystemTrayIcon *trayIcon;
    QMenu *trayMenu;
    void setupTrayIcon();
    void setupMenu();
};

#endif // TRAYICON_H
