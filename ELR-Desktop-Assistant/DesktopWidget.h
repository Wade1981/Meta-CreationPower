#ifndef DESKTOPWIDGET_H
#define DESKTOPWIDGET_H

#include <QWidget>
#include <QLabel>
#include <QVBoxLayout>
#include <QPushButton>
#include <QListWidget>

class DesktopWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DesktopWidget(QWidget *parent = nullptr);
    ~DesktopWidget();

public slots:
    void updateStatus(const QString &status);
    void updateContainers(const QList<QMap<QString, QString>> &containers);

protected:
    void mousePressEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;

private:
    QPoint lastMousePos;
    QLabel *statusLabel;
    QListWidget *containersList;
    QPushButton *refreshButton;
    QPushButton *hideButton;
    void setupUI();
};

#endif // DESKTOPWIDGET_H
