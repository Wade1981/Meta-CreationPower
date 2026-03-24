#ifndef ELRCLIENT_H
#define ELRCLIENT_H

#include <QObject>
#include <QNetworkAccessManager>
#include <QTimer>

class ELRClient : public QObject
{
    Q_OBJECT

public:
    explicit ELRClient(QObject *parent = nullptr);
    ~ELRClient();

public slots:
    void startMonitoring();
    void stopMonitoring();

signals:
    void statusUpdated(const QString &status);
    void containersUpdated(const QList<QMap<QString, QString>> &containers);

private slots:
    void checkELRStatus();
    void onStatusReplyFinished();
    void onContainersReplyFinished();

private:
    QNetworkAccessManager *networkManager;
    QTimer *monitoringTimer;
    QString elrApiBaseUrl;
    void getELRStatus();
    void getELRContainers();
};

#endif // ELRCLIENT_H
